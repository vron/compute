.text

.globl amd64_win_switch
.intel_syntax noprefix
/*
    https://en.wikipedia.org/wiki/X86_calling_conventions#x86-64_calling_conventions
    https://docs.microsoft.com/en-us/cpp/build/x64-calling-convention?view=vs-2019
    https://docs.microsoft.com/en-us/cpp/build/stack-usage?view=vs-2019
*/
amd64_win_switch:
    /* Store the current context */
    mov   QWORD PTR [rdx+0x0],       rbx
    mov   QWORD PTR [rdx+0x8],       rbp 
    mov   QWORD PTR [rdx+0x10],      rdi  
    mov   QWORD PTR [rdx+0x18],      rsi

    mov   r10,           QWORD PTR [rsp]    // return pos at top of stack
    mov   QWORD PTR [rdx+0x48],      r10    // return pos at top of stack
    lea   r10,                 [rsp+0x8]    // stack without return
    mov   QWORD PTR [rdx+0x20],      r10    // stack without return

    mov   QWORD PTR [rdx+0x28],      r12  
    mov   QWORD PTR [rdx+0x30],      r13 
    mov   QWORD PTR [rdx+0x38],      r14  
    mov   QWORD PTR [rdx+0x40],      r15

    movaps XMMWORD PTR [rdx+0x50],  xmm6    // calling con. also uses sse regs
    movaps XMMWORD PTR [rdx+0x60],  xmm7
    movaps XMMWORD PTR [rdx+0x70],  xmm8
    movaps XMMWORD PTR [rdx+0x80],  xmm9
    movaps XMMWORD PTR [rdx+0x90], xmm10
    movaps XMMWORD PTR [rdx+0xA0], xmm11
    movaps XMMWORD PTR [rdx+0xB0], xmm12
    movaps XMMWORD PTR [rdx+0xC0], xmm13
    movaps XMMWORD PTR [rdx+0xD0], xmm14
    movaps XMMWORD PTR [rdx+0xE0], xmm15

    /* restore the destination context */
    mov   rbx,       QWORD PTR [rcx+0x0]
    mov   rbp,       QWORD PTR [rcx+0x8]
    mov   rdi,      QWORD PTR [rcx+0x10]   
    mov   rsi,      QWORD PTR [rcx+0x18]

    mov   rsp,      QWORD PTR [rcx+0x20]

    mov   r12,      QWORD PTR [rcx+0x28]
    mov   r13,      QWORD PTR [rcx+0x30]
    mov   r14,      QWORD PTR [rcx+0x38]
    mov   r15,      QWORD PTR [rcx+0x40]

    movaps xmm6,  XMMWORD PTR [rcx+0x50]
    movaps xmm7,  XMMWORD PTR [rcx+0x60] 
    movaps xmm8,  XMMWORD PTR [rcx+0x70]
    movaps xmm9,  XMMWORD PTR [rcx+0x80] 
    movaps xmm10, XMMWORD PTR [rcx+0x90] 
    movaps xmm11, XMMWORD PTR [rcx+0xA0] 
    movaps xmm12, XMMWORD PTR [rcx+0xB0] 
    movaps xmm13, XMMWORD PTR [rcx+0xC0] 
    movaps xmm14, XMMWORD PTR [rcx+0xD0] 
    movaps xmm15, XMMWORD PTR [rcx+0xE0]

    mov   r10,      QWORD PTR [rcx+0x48]
    jmp   r10
