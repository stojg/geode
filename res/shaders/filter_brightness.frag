#version 410 core

out vec3 FragColor;

in vec2 TexCoords;

uniform sampler2D x_filterTexture;

float min = 0.8;
float max = 2.0;

void main()
{
    vec3 texCol = texture(x_filterTexture, TexCoords).rgb;
    vec3 luminanceVector = vec3(0.2125, 0.7154, 0.0721);
    float luminance = dot(texCol.rgb, luminanceVector);
    FragColor = texCol * 4.0 * smoothstep(min, max, luminance);
    FragColor = clamp(FragColor, 0.0, 16);
}
