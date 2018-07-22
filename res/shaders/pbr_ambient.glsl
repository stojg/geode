uniform sampler2D   x_brdfLUT;
uniform samplerCube x_irradianceMap;
uniform samplerCube x_prefilterMap;
vec3 CalcAmbient(vec3 normal,vec3 V, vec3 F0, Mtrl mtrl)
{
    vec3 F = fresnelSchlickRoughness(max(dot(normal, V), 0.0), F0, mtrl.roughness);

    vec3 kS = F;
    vec3 kD = 1.0 - kS;
    kD *= 1.0 - mtrl.metallic;

    // diffuse
    vec3 irradiance = texture(x_irradianceMap, vs_in.Normal).rgb;
    vec3 diffuse    = irradiance * mtrl.albedo;

    // specular
    const float MAX_REFLECTION_LOD = 4.0;
    vec3 prefilteredColor = textureLod(x_prefilterMap, vs_in.Reflection,  mtrl.roughness * MAX_REFLECTION_LOD).rgb;
    vec2 brdf  = texture(x_brdfLUT, vec2(max(dot(normal, V), 0.0), mtrl.roughness)).rg;
    vec3 specular = prefilteredColor * (F * brdf.x + brdf.y);

    // sum up all ambient
    return (kD * diffuse + specular);
}


vec3 CalcAmbient(vec3 normal,vec3 V, vec3 F0, Mtrl mtrl, vec3 R)
{
    vec3 F = fresnelSchlickRoughness(max(dot(normal, V), 0.0), F0, mtrl.roughness);

    vec3 kS = F;
    vec3 kD = 1.0 - kS;
    kD *= 1.0 - mtrl.metallic;

    // diffuse
    vec3 irradiance = texture(x_irradianceMap, vs_in.Normal).rgb;
    vec3 diffuse    = irradiance * mtrl.albedo;

    // specular
    const float MAX_REFLECTION_LOD = 4.0;
    vec3 prefilteredColor = textureLod(x_prefilterMap, R,  mtrl.roughness * MAX_REFLECTION_LOD).rgb;
    vec2 brdf  = texture(x_brdfLUT, vec2(max(dot(normal, V), 0.0), mtrl.roughness)).rg;
    vec3 specular = prefilteredColor * (F * brdf.x + brdf.y);

    // sum up all ambient
    return (kD * diffuse + specular);
}
