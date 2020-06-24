# Test of using GLSL compute kernals on CPU in go

Not working at the moment. To try it run the following in the main folder.

    (cd build && docker build . -t temp && (cd .. && docker run -v $(pwd):/data temp))