#version 410

uniform vec2 screenSize;

in vec2 position;

in mat2x3 dstTransform;
in float dstDepth;
in vec3 srcLocation;
in vec2 srcSize;

out vec3 f_srcLocation;

void main() {
	mat4 m;
    m[0] = vec4(dstTransform[0][0],	dstTransform[1][0],	0.0,	0.0);
    m[1] = vec4(dstTransform[0][1], dstTransform[1][1],	0.0,	0.0);
    m[2] = vec4(0.0,           		0.0,				1.0,	0.0);
    m[3] = vec4(dstTransform[0][2], dstTransform[1][2],	0.0,	1.0);

    gl_Position = m*vec4(srcSize*position, dstDepth, 1.0);
    gl_Position.xy /= screenSize/2.0;
    gl_Position.xy -= vec2(1.0, 1.0);
    f_srcLocation = srcLocation + vec3(srcSize*((position*vec2(1.0,-1.0))+vec2(0.0,1.0)), 0.0);
}