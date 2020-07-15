// This file will be generated as part of the build, it will
// contain the body used to set members in the c++ shader struct
// from the parameters passed through the c - api

	if((((uintptr_t)(const void *)(d.din)) % (8)) != 0) { return this->set_error(EINVAL, "the argument d.din provided was not aligned to a 8 byte boundary as required"); };
	this->din = (float*)d.din;
	if((((uintptr_t)(const void *)(d.dout)) % (8)) != 0) { return this->set_error(EINVAL, "the argument d.dout provided was not aligned to a 8 byte boundary as required"); };
	this->dout = (float*)d.dout;
	return 0;