package pickle

import (
	"reflect"
	"fmt"
	"context"
)

func (cfg *frozenConfig) addDecoderToCache(cacheKey reflect.Type, decoder RootDecoder) {
	cfg.decoderCache.Store(cacheKey, decoder)
}

func (cfg *frozenConfig) addEncoderToCache(cacheKey reflect.Type, encoder RootEncoder) {
	cfg.encoderCache.Store(cacheKey, encoder)
}

func (cfg *frozenConfig) getDecoderFromCache(cacheKey reflect.Type) RootDecoder {
	decoder, found := cfg.decoderCache.Load(cacheKey)
	if found {
		return decoder.(RootDecoder)
	}
	return nil
}

func (cfg *frozenConfig) getEncoderFromCache(cacheKey reflect.Type) RootEncoder {
	encoder, found := cfg.encoderCache.Load(cacheKey)
	if found {
		return encoder.(RootEncoder)
	}
	return nil
}

func encoderOfType(cfg *frozenConfig, valType reflect.Type) (RootEncoder, error) {
	cacheKey := valType
	rootEncoder := cfg.getEncoderFromCache(cacheKey)
	if rootEncoder != nil {
		return rootEncoder, nil
	}
	encoder, err := createEncoderOfType(cfg, valType)
	if err != nil {
		return nil, err
	}
	rootEncoder = wrapRootEncoder(encoder)
	cfg.addEncoderToCache(cacheKey, rootEncoder)
	return rootEncoder, err
}

func wrapRootEncoder(encoder ValEncoder) RootEncoder {
	valType := encoder.Type()
	valKind := valType.Kind()
	rootEncoder := rootEncoder{valType, encoder.Signature(), encoder}
	switch valKind {
	case reflect.Struct:
		if valType.NumField() == 1 && valType.Field(0).Type.Kind() == reflect.Ptr {
			return &singlePointerFix{rootEncoder}
		}
	case reflect.Array:
		if valType.Len() == 1 && valType.Elem().Kind() == reflect.Ptr {
			return &singlePointerFix{rootEncoder}
		}
	case reflect.Ptr:
		return &singlePointerFix{rootEncoder}
	}
	return &rootEncoder
}

func decoderOfType(cfg *frozenConfig, valType reflect.Type) (RootDecoder, error) {
	cacheKey := valType
	rootDecoder := cfg.getDecoderFromCache(cacheKey)
	if rootDecoder != nil {
		return rootDecoder, nil
	}
	decoder, err := createDecoderOfType(cfg, valType)
	if err != nil {
		return nil, err
	}
	if cfg.readonlyDecode && decoder.HasPointer() {
		rootDecoder = &rootDecoderWithCopy{valType, decoder.Signature(), decoder}
	} else {
		rootDecoder = &rootDecoderWithoutCopy{valType, decoder.Signature(), decoder}
	}
	cfg.addDecoderToCache(cacheKey, rootDecoder)
	return rootDecoder, err
}

var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()

func createEncoderOfType(cfg *frozenConfig, valType reflect.Type) (ValEncoder, error) {
	if valType == contextType {
		return &NoopCodec{BaseCodec: *newBaseCodec(valType, uint64(valType.Kind()))}, nil
	}
	valKind := valType.Kind()
	switch valKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Float32, reflect.Float64:
		return &NoopCodec{BaseCodec: *newBaseCodec(valType, uint64(valKind))}, nil
	case reflect.String:
		return &stringCodec{BaseCodec: *newBaseCodec(valType, uint64(valKind))}, nil
	case reflect.Interface:
		return &interfaceEncoder{BaseCodec: *newBaseCodec(valType, uint64(valKind))}, nil
	case reflect.Struct:
		signature := uint64(valKind)
		fields := make([]structFieldEncoder, 0, valType.NumField())
		for i := 0; i < valType.NumField(); i++ {
			encoder, err := createEncoderOfType(cfg, valType.Field(i).Type)
			if err != nil {
				return nil, err
			}
			signature = 31*signature + encoder.Signature()
			if !encoder.IsNoop() {
				fields = append(fields, structFieldEncoder{
					offset:  valType.Field(i).Offset,
					encoder: encoder,
				})
			}
		}
		encoder := &structEncoder{BaseCodec: *newBaseCodec(valType, signature), fields: fields}
		return encoder, nil
	case reflect.Array:
		signature := uint64(valKind)
		elemEncoder, err := createEncoderOfType(cfg, valType.Elem())
		if err != nil {
			return nil, err
		}
		signature = 31*signature + elemEncoder.Signature()
		if elemEncoder.IsNoop() {
			elemEncoder = nil
		}
		encoder := &arrayEncoder{
			BaseCodec:   *newBaseCodec(valType, signature),
			arrayLength: valType.Len(),
			elementSize: valType.Elem().Size(),
			elemEncoder: elemEncoder,
		}
		return encoder, nil
	case reflect.Slice:
		signature := uint64(valKind)
		elemEncoder, err := createEncoderOfType(cfg, valType.Elem())
		if err != nil {
			return nil, err
		}
		signature = 31*signature + elemEncoder.Signature()
		if elemEncoder.IsNoop() {
			elemEncoder = nil
		}
		return &sliceEncoder{BaseCodec: *newBaseCodec(valType, signature),
			elemSize: int(valType.Elem().Size()), elemEncoder: elemEncoder}, nil
	case reflect.Ptr:
		signature := uint64(valKind)
		elemEncoder, err := createEncoderOfType(cfg, valType.Elem())
		if err != nil {
			return nil, err
		}
		signature = 31*signature + elemEncoder.Signature()
		encoder := &pointerEncoder{BaseCodec: *newBaseCodec(valType, signature), elemEncoder: elemEncoder}
		return encoder, nil
	}
	return nil, fmt.Errorf("unsupported type %s", valType.String())
}

func createDecoderOfType(cfg *frozenConfig, valType reflect.Type) (ValDecoder, error) {
	valKind := valType.Kind()
	switch valKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Float32, reflect.Float64:
		return &NoopCodec{BaseCodec: *newBaseCodec(valType, uint64(valKind))}, nil
	case reflect.String:
		return &stringCodec{BaseCodec: *newBaseCodec(valType, uint64(valKind))}, nil
	case reflect.Struct:
		fields := make([]structFieldDecoder, 0, valType.NumField())
		signature := uint64(valKind)
		hasPointer := false
		for i := 0; i < valType.NumField(); i++ {
			decoder, err := createDecoderOfType(cfg, valType.Field(i).Type)
			if err != nil {
				return nil, err
			}
			if decoder.HasPointer() {
				hasPointer = true
			}
			signature = 31*signature + decoder.Signature()
			if !decoder.IsNoop() {
				fields = append(fields, structFieldDecoder{
					offset:  valType.Field(i).Offset,
					decoder: decoder,
				})
			}
		}
		if hasPointer {
			return &structDecoderWithPointer{BaseCodec: *newBaseCodec(valType, signature), fields: fields}, nil
		}
		return &structDecoderWithoutPointer{BaseCodec: *newBaseCodec(valType, signature), fields: fields}, nil
	case reflect.Array:
		signature := uint64(valKind)
		elemDecoder, err := createDecoderOfType(cfg, valType.Elem())
		if err != nil {
			return nil, err
		}
		signature = 31*signature + elemDecoder.Signature()
		hasPointer := elemDecoder.HasPointer()
		if elemDecoder.IsNoop() {
			elemDecoder = nil
		}
		if hasPointer {
			return &arrayDecoderWithPointer{
				BaseCodec:   *newBaseCodec(valType, signature),
				arrayLength: valType.Len(),
				elementSize: valType.Elem().Size(),
				elemDecoder: elemDecoder,
			}, nil
		}
		return &arrayDecoderWithoutPointer{
			BaseCodec:   *newBaseCodec(valType, signature),
			arrayLength: valType.Len(),
			elementSize: valType.Elem().Size(),
			elemDecoder: elemDecoder,
		}, nil
	case reflect.Slice:
		signature := uint64(valKind)
		elemDecoder, err := createDecoderOfType(cfg, valType.Elem())
		if err != nil {
			return nil, err
		}
		signature = 31*signature + elemDecoder.Signature()
		shouldCopy := false
		if elemDecoder.HasPointer() && cfg.readonlyDecode {
			shouldCopy = true
		}
		if elemDecoder.IsNoop() {
			elemDecoder = nil
		}
		if shouldCopy {
			return &sliceDecoderWithCopy{BaseCodec: *newBaseCodec(valType, signature),
				elemSize: int(valType.Elem().Size()), elemDecoder: elemDecoder}, nil
		}
		return &sliceDecoderWithoutCopy{BaseCodec: *newBaseCodec(valType, signature),
			elemSize: int(valType.Elem().Size()), elemDecoder: elemDecoder}, nil
	case reflect.Ptr:
		signature := uint64(valKind)
		elemDecoder, err := createDecoderOfType(cfg, valType.Elem())
		if err != nil {
			return nil, err
		}
		signature = 31*signature + elemDecoder.Signature()
		if elemDecoder.HasPointer() && cfg.readonlyDecode {
			return &pointerDecoderWithCopy{BaseCodec: *newBaseCodec(valType, signature), elemDecoder: elemDecoder}, nil
		}
		return &pointerDecoderWithoutCopy{BaseCodec: *newBaseCodec(valType, signature), elemDecoder: elemDecoder}, nil
	}
	return nil, fmt.Errorf("unsupported type %s", valType.String())
}
