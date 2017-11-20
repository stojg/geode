#version 410

out vec4 fragColor;

uniform sampler2D diffuse;

in vec2 TexCoord;

void main() {
    fragColor = texture(diffuse, TexCoord);
}
