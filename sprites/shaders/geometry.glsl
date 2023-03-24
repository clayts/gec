#version 410
layout (points) in;
layout (triangle_strip, max_vertices = 4) out;

in mat2x3 g_dstTransform[];
in float g_dstDepth[];
in vec3 g_srcLocation[];
in vec2 g_srcSize[];
out vec3 f_srcLocation;

uniform vec2 screenSize;

void main() {
	mat4 m;
    m[0] = vec4(g_dstTransform[0][0][0],	g_dstTransform[0][1][0],	0.0,	0.0);
    m[1] = vec4(g_dstTransform[0][0][1],    g_dstTransform[0][1][1],	0.0,	0.0);
    m[2] = vec4(0.0,                        0.0,						1.0,	0.0);
    m[3] = vec4(g_dstTransform[0][0][2],    g_dstTransform[0][1][2],	0.0,	1.0);

    gl_Position = m*vec4(0.0, 0.0, g_dstDepth[0], 1.0);
    gl_Position.xy /= screenSize/2.0;
    gl_Position.xy -= vec2(1.0, 1.0);
    f_srcLocation = g_srcLocation[0] + vec3(g_srcSize[0]*vec2(0.0, 1.0), 0.0);
    EmitVertex();

    gl_Position = m*vec4(0.0, g_srcSize[0].y, g_dstDepth[0], 1.0);
    gl_Position.xy /= screenSize/2.0;
    gl_Position.xy -= vec2(1.0, 1.0);
    f_srcLocation = g_srcLocation[0] + vec3(g_srcSize[0]*vec2(0.0, 0.0), 0.0);
    EmitVertex();
        
    gl_Position = m*vec4(g_srcSize[0].x, 0.0, g_dstDepth[0], 1.0);
    gl_Position.xy /= screenSize/2.0;
    gl_Position.xy -= vec2(1.0, 1.0);
    f_srcLocation = g_srcLocation[0] + vec3(g_srcSize[0]*vec2(1.0, 1.0), 0.0);
    EmitVertex();
    
    gl_Position = m*vec4(g_srcSize[0].x, g_srcSize[0].y, g_dstDepth[0], 1.0);
    gl_Position.xy /= screenSize/2.0;
    gl_Position.xy -= vec2(1.0, 1.0);
    f_srcLocation = g_srcLocation[0] + vec3(g_srcSize[0]*vec2(1.0, 0.0), 0.0);
    EmitVertex();

    EndPrimitive();
}  