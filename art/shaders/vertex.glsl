#version 410

uniform vec2 screenSize;
uniform mat2x3 cameraTransform;

in vec2 position;

in mat2x3 dstTransform;
in float dstDepth;
in vec4 srcLocation;
in vec2 srcSize;

out vec4 f_srcLocation;

void main() {
	mat4 m;
    m[0] = vec4(dstTransform[0][0],	dstTransform[1][0],	0.0,	0.0);
    m[1] = vec4(dstTransform[0][1], dstTransform[1][1],	0.0,	0.0);
    m[2] = vec4(0.0,           		0.0,				1.0,	0.0);
    m[3] = vec4(dstTransform[0][2], dstTransform[1][2],	0.0,	1.0);

	mat4 c;
    c[0] = vec4(cameraTransform[0][0], 	cameraTransform[1][0],	0.0,	0.0);
    c[1] = vec4(cameraTransform[0][1], 	cameraTransform[1][1],	0.0,	0.0);
    c[2] = vec4(0.0,           			0.0,					1.0,	0.0);
    c[3] = vec4(cameraTransform[0][2], 	cameraTransform[1][2],	0.0,	1.0);

    gl_Position = c*m*vec4(srcSize*position, dstDepth, 1.0);
    gl_Position.xy /= screenSize/2.0;
    gl_Position.xy -= vec2(1.0, 1.0);
    f_srcLocation = srcLocation;
    f_srcLocation.xy += srcSize*position;
}