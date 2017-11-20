#version 410 core

uniform sampler2D diffuse;

in vec2 TexCoord;
out vec4 fragColor;

const float ambientStrength = 0.01;

void main() {

    fragColor = texture(diffuse, TexCoord) * ambientStrength;
}
