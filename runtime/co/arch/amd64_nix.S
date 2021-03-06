.text

.globl amd64_nix_switch
#ifndef __APPLE__
.type  amd64_nix_switch, @function
#endif
.intel_syntax noprefix
/*
    TODO: Do we need the ctrl words? (fnstcw,stmxcsr,fldcwldmxcsr) 
    https://en.wikipedia.org/wiki/X86_calling_conventions
    https://aaronbloomfield.github.io/pdr/book/x86-64bit-ccc-chapter.pdf
*/
amd64_nix_switch:
    /* Store the current context */
    mov   QWORD PTR [rsi+0x0], rbx    // register assumed unm. through a call
    mov   QWORD PTR [rsi+0x8], rbp    // register assumed unm. through a call
    mov   QWORD PTR [rsi+0x10],r12    // register assumed unm. through a call
    mov   QWORD PTR [rsi+0x18],r13    // register assumed unm. through a call
    mov   QWORD PTR [rsi+0x20],r14    // register assumed unm. through a call
    mov   QWORD PTR [rsi+0x28],r15    // register assumed unm. through a call
    mov   rdx,     QWORD PTR [rsp]    // load the return address from top of stack
    mov   QWORD PTR [rsi+0x30],rdx
    lea   rcx,           [rsp+0x8]    // stack without return, lea to not change any flgs
    mov   QWORD PTR [rsi+0x38],rcx

    /* restore the destination context */
    mov   rbx, QWORD PTR [rdi+0x0]
    mov   rbp, QWORD PTR [rdi+0x8]
    mov   r12,QWORD PTR [rdi+0x10]
    mov   r13,QWORD PTR [rdi+0x18]
    mov   r14,QWORD PTR [rdi+0x20]
    mov   r15,QWORD PTR [rdi+0x28]
    mov   rax,QWORD PTR [rdi+0x30]
    mov   rcx,QWORD PTR [rdi+0x38]
    mov   rsp,rcx
    jmp   rax
