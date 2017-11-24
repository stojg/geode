#version 410 core

out vec4 FragColor;

in vec2 TexCoords;

uniform sampler2D x_filterTexture;
uniform vec3 x_blurScale;

void main()
{

    vec4 color = vec4(0.0);

    vec2 scale = x_blurScale.rg;

    color += texture(x_filterTexture, TexCoords.xy + (vec2(-3.0) * scale)) * vec4(0.015625);
    color += texture(x_filterTexture, TexCoords.xy + (vec2(-2.0) * scale)) * vec4(0.09375);
    color += texture(x_filterTexture, TexCoords.xy + (vec2(-1.0) * scale)) * vec4(0.234375);
    color += texture(x_filterTexture, TexCoords.xy + (vec2(+0.0) * scale)) * vec4(0.3125);
    color += texture(x_filterTexture, TexCoords.xy + (vec2(+1.0) * scale)) * vec4(0.234375);
    color += texture(x_filterTexture, TexCoords.xy + (vec2(+2.0) * scale)) * vec4(0.09375);
    color += texture(x_filterTexture, TexCoords.xy + (vec2(+3.0) * scale)) * vec4(0.015625);

    FragColor = color;
}
