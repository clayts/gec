#version 410
in mat2x3 dstTransform;
in float dstDepth;
in vec3 srcLocation;
in vec2 srcSize;
out mat2x3 g_dstTransform;
out float g_dstDepth;
out vec3 g_srcLocation;
out vec2 g_srcSize;
void main() {
	g_dstTransform = dstTransform;
	g_srcLocation = srcLocation;
	g_srcSize = srcSize;
	g_dstDepth = dstDepth;
}