#version 410 core

out vec4 FragColor;

in vec2 TexCoords;

uniform sampler2D x_filterTexture;
uniform vec3 x_blurScale;

uniform float weight[3] = float[] (0.2270270270, 0.3162162162, 0.0702702703 );
uniform float offset[3] = float[]( 0.0, 1.3846153846, 3.2307692308 );

void main()
{

    vec2 scale = x_blurScale.rg;

    //vec2 tex_offset = 1.0 / textureSize(x_filterTexture, 0); // gets size of single texel
    //vec4 result = vec4(texture(x_filterTexture, TexCoords).rgb * weight[0], 0);
    //for(int i = 1; i<3; i++) {
      //  result += texture(x_filterTexture, TexCoords + vec2(tex_offset.x * offset[i]) * scale) * vec4(weight[i]);
     //   result += texture(x_filterTexture, TexCoords - vec2(tex_offset.x * offset[i]) * scale) * vec4(weight[i]);
    //}
    //FragColor = result;
    //return;

    vec4 color = vec4(0.0);

    color += texture(x_filterTexture, TexCoords + (vec2(-3.0) * scale)) * vec4(0.015625);
    color += texture(x_filterTexture, TexCoords + (vec2(-2.0) * scale)) * vec4(0.09375);
    color += texture(x_filterTexture, TexCoords + (vec2(-1.0) * scale)) * vec4(0.234375);
    color += texture(x_filterTexture, TexCoords + (vec2(+0.0) * scale)) * vec4(0.3125);
    color += texture(x_filterTexture, TexCoords + (vec2(+1.0) * scale)) * vec4(0.234375);
    color += texture(x_filterTexture, TexCoords + (vec2(+2.0) * scale)) * vec4(0.09375);
    color += texture(x_filterTexture, TexCoords + (vec2(+3.0) * scale)) * vec4(0.015625);

    FragColor = color;
}
