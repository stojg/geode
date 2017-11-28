#version 410 core

out vec4 FragColor;

in vec2 TexCoords;

uniform sampler2D x_filterTexture;

void main()
{
    vec4 color = texture(x_filterTexture, TexCoords).rgba;
    if (color.z < 0.0001) {
        discard;
    }
    FragColor = vec4(color.rgb, 1.0);
}
