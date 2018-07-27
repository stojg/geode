#version 410 core

out vec4 FragColor;

in vec2 TexCoords;

uniform sampler2D x_filterTexture;
uniform sampler2D x_filterTexture2;
uniform sampler2D x_filterTexture3;
uniform sampler2D x_filterTexture4;

void main()
{

    vec3 hdrColor1 = texture(x_filterTexture, TexCoords).rgb;
    vec3 hdrColor2 = texture(x_filterTexture2, TexCoords).rgb;
    vec3 hdrColor3 = texture(x_filterTexture3, TexCoords).rgb;
    vec3 hdrColor4 = texture(x_filterTexture4, TexCoords).rgb;
    FragColor = vec4(vec3(hdrColor1 + (hdrColor2 + hdrColor3 + hdrColor4) * 0.3), 1);

}
