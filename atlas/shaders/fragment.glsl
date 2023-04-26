#version 410

uniform sampler2DArray textureArray[16];

in vec4 f_srcLocation;

out vec4 frag_colour;

void main() {
    frag_colour = texelFetch(textureArray[int(f_srcLocation.w)], ivec3(int(f_srcLocation.x), int(f_srcLocation.y), int(f_srcLocation.z)), 0);
    if (frag_colour.a == 0.0) {
        discard;
    }
}