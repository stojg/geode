#version 410

uniform sampler2D diffuse;

in vec2 TexCoord;
out vec4 fragColor;

void main() {
    float ambientStrength = 0.01;
    fragColor = texture(diffuse, TexCoord) * ambientStrength;
}
