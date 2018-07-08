#version 410 core

out vec4 FragColor;

in vec2 TexCoords;

uniform sampler2D x_filterTexture;

void main()
{
    float Pixels = 1024.0*2;
    float dx = 15.0 * (1.0 / Pixels);
    float dy = 10.0 * (1.0 / Pixels);
    vec2 Coord = vec2(dx * floor(TexCoords.x / dx), dy * floor(TexCoords.y / dy));
    FragColor = texture(x_filterTexture, Coord);
}
