// Copyright 2016 Huan Du. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "textflag.h"

TEXT ·getg(SB), NOSPLIT, $0-8
    MOV    g, A0
    MOV    A0, ret+0(FP)
    RET
