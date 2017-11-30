#version 410 core

out vec4 FragColor;

uniform sampler2D diffuse;
uniform vec3 x_lightColors[16];
uniform float specularStrength = 0.1;
uniform int x_numPointLights;

in vec2 TexCoord;
in vec3 Normal;
in vec3 FragPos;
in vec3 LightPositions[16];
in vec3 ModelViewPos;

vec3 CalcPointLight(vec3 lightPosition, vec3 lightColor, vec3 objectColor, vec3 norm) {


    vec3 lightDiff = lightPosition - ModelViewPos;
    float distance = length(lightDiff);

    if (distance > 30) {
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
    vec3 halfwayDir = normalize(lightDirection - normalize(ModelViewPos));
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
    vec3 objectColor = texture(diffuse, TexCoord).rgb;
    vec3 norm = normalize(Normal);

    vec3 final = vec3(0);

    float ambientStrength = 0.01;
    final += ambientStrength * objectColor;

    for (int i = 0; i < x_numPointLights; i++) {
        final += CalcPointLight(LightPositions[i], x_lightColors[i], objectColor, norm);
    }

    FragColor = vec4(final, 0);
}
