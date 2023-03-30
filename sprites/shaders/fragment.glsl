#version 410

uniform sampler2DArray textureArray;

in vec3 f_srcLocation;

out vec4 frag_colour;

void main() {
    frag_colour = texelFetch(textureArray, ivec3(int(f_srcLocation.x), int(f_srcLocation.y), int(f_srcLocation.z)), 0);
    if (frag_colour.a == 0.0) {
        discard;
    }
}