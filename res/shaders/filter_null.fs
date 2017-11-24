#version 410 core

out vec4 FragColor;

in vec2 TexCoords;

uniform sampler2D x_filterTexture;

void main()
{
    vec3 color = texture(x_filterTexture, TexCoords).rgb;
    FragColor = vec4(vec3(1), 1.0);
}
