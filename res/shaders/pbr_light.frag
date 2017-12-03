#version 410 core

struct Material {
    vec3 albedo;
    float metallic;
    float roughness;
};

uniform sampler2D   x_brdfLUT;
uniform samplerCube x_irradianceMap;
uniform samplerCube x_prefilterMap;

uniform mat4 InverseMV;
uniform mat4 InvView;
uniform mat4 view;

#include "pbr_lights.glsl"
#include "pbr.glsl"

void main() {

    // Normal in view space
    vec3 normal = normalize(vs_in.V_Normal);

    vec3 F0 = vec3(0.04);
    F0 = mix(F0, material.albedo, material.metallic);

    vec3 Lo = vec3(0.0);

    vec3 V = normalize(-vs_in.W_ViewPos);

    for (int i = 0; i < numLights; i++) {
        if (lights[i].constant == 0) {
            Lo += CalcDirectional(F0, vs_in.V_LightPositions[i], lights[i], material, normal, vs_in.W_ViewPos, V);
        } else if (lights[i].cutoff > 0) {
            Lo += CalcSpot(F0, vs_in.V_LightPositions[i], lights[i], material, normal, vs_in.W_ViewPos, V);
        } else {
            Lo += CalcPoint(F0, vs_in.V_LightPositions[i], lights[i], material, normal, vs_in.W_ViewPos, V);
        }
    }

    if (x_enable_env_map == 0) {
        FragColor = vec4(Lo, 1);
        return;
    }
    // enviroment ambient lightning

    // direction towards they eye (camera) in the view (eye) space
    vec3 viewDirection = normalize(-vs_in.W_ViewPos);
    // eye direction in worldspace
    vec3 wcEyeDir = vec3(InvView * vec4(viewDirection, 0.0));

    // reflection
    vec3 R = reflect(-wcEyeDir, normalize(vs_in.Normal));

    vec3 F = fresnelSchlickRoughness(max(dot(vs_in.V_Normal, viewDirection), 0.0), F0, material.roughness);

    vec3 kS = F;
    vec3 kD = 1.0 - kS;
    kD *= 1.0 - material.metallic;

    // diffuse
    vec3 irradiance = texture(x_irradianceMap, vs_in.Normal).rgb;
    vec3 diffuse    = irradiance * material.albedo;

    // specular
    const float MAX_REFLECTION_LOD = 4.0;
    vec3 prefilteredColor = textureLod(x_prefilterMap, R,  material.roughness * MAX_REFLECTION_LOD).rgb;
    vec2 brdf  = texture(x_brdfLUT, vec2(max(dot(vs_in.Normal, wcEyeDir), 0.0), material.roughness)).rg;
    vec3 specular = prefilteredColor * (F * brdf.x + brdf.y);

    // sum up all ambient
    vec3 ambient = (kD * diffuse + specular);

    // combine with lights
    vec3 color = Lo + ambient;

    FragColor = vec4(color, 1);
}
