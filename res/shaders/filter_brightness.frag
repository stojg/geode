#version 410 core

out vec4 FragColor;

in vec2 TexCoords;

uniform sampler2D x_filterTexture;

void main()
{
    vec3 luminanceVector = vec3(0.2125, 0.7154, 0.0721);
    vec3 texCol = texture(x_filterTexture, TexCoords).rgb;
    float luminance = dot(luminanceVector, texCol.rgb);
    luminance = (atan((luminance-1.9)*1024) / 3.141592) + 0.5;
    FragColor = vec4(texCol.rgb * luminance, 1.0);
    FragColor = clamp(FragColor, 0.0, 2.0);
}
