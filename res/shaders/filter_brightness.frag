#version 410 core

out vec4 FragColor;

in vec2 TexCoords;

uniform sampler2D x_filterTexture;

void main()
{

    vec3 luminanceVector = vec3(0.2125, 0.7154, 0.0721);
    vec3 texCol = texture(x_filterTexture, TexCoords).rgb;
    float luminance = dot(luminanceVector, texCol.rgb);
    FragColor = vec4(texCol.rgb * luminance * luminance * luminance, 1.0);
    FragColor = clamp(FragColor, 0, 5);
//    if (luminance > 1.2) {
//        FragColor = vec4(texCol.rgb, 1.0);
//    } else {
//        FragColor = vec4(0, 0, 0, 1.0);
//    }
}
