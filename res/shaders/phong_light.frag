#version 410 core

in VS_OUT
{
    vec3 V_Normal;
    vec2 TexCoord;
    vec3 V_LightPositions[16];
    vec3 W_ViewPos;
} vs_in;

uniform sampler2D diffuse;
uniform float specularStrength = 0.1;


struct Light {
    vec3 position;
    vec3 color;
};
uniform int numPointLights;
uniform Light pointLights[16];

out vec4 FragColor;

vec3 CalcPointLight(vec3 lightPosition, vec3 lightColor, vec3 objectColor, vec3 norm, vec3 viewPos) {
    vec3 lightDiff = lightPosition - viewPos;
    float distance = length(lightDiff);

    if (distance > 20) {
        return vec3(0);
    }
    vec3 lightDirection = normalize(lightDiff);

    const float constant = 1.0;
    const float linear = 0.7;
    const float quadratic = 1.8;
    float attenuation = 1.0 / (constant + linear * distance + quadratic * (distance * distance));

    // diffuse
    float diff = max(dot(norm, lightDirection), 0.0);

    // specular
    vec3 halfwayDir = normalize(lightDirection - normalize(viewPos));
    vec3 reflectDir = reflect(-lightDirection, norm);
    float spec = pow(max(dot(norm, halfwayDir), 0.0), 16);

    // combine results
    vec3 diffuseColor = lightColor * diff * objectColor;
    vec3 specularColor = lightColor * spec * objectColor;

    diffuseColor *= attenuation;
    specularColor *= attenuation;

    return diffuseColor + specularColor;
}

void main() {
    vec3 objectColor = texture(diffuse, vs_in.TexCoord).rgb;
    vec3 normal = normalize(vs_in.V_Normal);

    vec3 final = vec3(0);

    float ambientStrength = 0.01;
    final += ambientStrength * objectColor;

    for (int i = 0; i < numPointLights; i++) {
        final += CalcPointLight(vs_in.V_LightPositions[i], pointLights[i].color, objectColor, normal, vs_in.W_ViewPos);
    }

    FragColor = vec4(final, 0);
}
