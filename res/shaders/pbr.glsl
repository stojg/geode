
const float PI = 3.14159265359;

struct Mtrl {
    vec3 albedo;
    float metallic;
    float roughness;
};

// The Fresnel-Schlick approximation expects a F0 parameter which is known as the surface reflection at zero incidence
// or how much the surface reflects if looking directly at the surface, ie it calculates the ratio between specular and
// diffuse reflection, or how much the surface reflects light versus how much it refracts light.
vec3 fresnelSchlick(float cosTheta, vec3 F0)
{
    return F0 + (1.0 - F0) * pow(1.0 - cosTheta, 5.0);
}

vec3 fresnelSchlickRoughness(float cosTheta, vec3 F0, float roughness)
{
    return F0 + (max(vec3(1.0 - roughness), F0) - F0) * pow(1.0 - cosTheta, 5.0);
}

float DistributionGGX(vec3 N, vec3 H, float roughness)
{
    float a      = roughness*roughness;
    float a2     = a*a;
    float NdotH  = max(dot(N, H), 0.0);
    float NdotH2 = NdotH*NdotH;

    float nom   = a2;
    float denom = (NdotH2 * (a2 - 1.0) + 1.0);
    denom = PI * denom * denom;

    return nom / denom;
}

float GeometrySchlickGGX(float NdotV, float roughness)
{
    float r = (roughness + 1.0);
    float k = (r*r) / 8.0;

    float nom   = NdotV;
    float denom = NdotV * (1.0 - k) + k;

    return nom / denom;
}

float GeometrySmith(vec3 N, vec3 V, vec3 L, float roughness)
{
    float NdotV = max(dot(N, V), 0.0);
    float NdotL = max(dot(N, L), 0.0);
    float ggx2  = GeometrySchlickGGX(NdotV, roughness);
    float ggx1  = GeometrySchlickGGX(NdotL, roughness);

    return ggx1 * ggx2;
}

uniform Material material;

uniform int numLights;
uniform Light lights[16];

uniform int x_enable_env_map;

out vec4 FragColor;

vec3 calcCookTorrance(vec3 H, vec3 V, vec3 N, Mtrl material, vec3 F0, vec3 L, vec3 radiance) {
    vec3 F    = fresnelSchlick(max(dot(H, V), 0.0), F0);
    float NDF = DistributionGGX(N, H, material.roughness);
    float G   = GeometrySmith(N, V, L, material.roughness);

    vec3 nominator    = NDF * G * F;
    float denominator = 4.0 * max(dot(N, V), 0.0) * max(dot(N, L), 0.0);
    vec3 specular     = nominator / max(denominator, 0.001);

    vec3 kS = F;
    vec3 kD = vec3(1.0) - kS;
    kD *= 1.0 - material.metallic;
    float NdotL = max(dot(N, L), 0.0);
    return (kD * material.albedo / PI + specular) * radiance * NdotL;
}

vec3 CalcPoint(vec3 F0, vec3 lightPosition, Light light, Mtrl material,  vec3 N, vec3 viewPos, vec3 V) {

    vec3 viewLightDirection = (view * vec4(light.direction, 0)).xyz;
    float dist = distance(lightPosition, viewPos);
    if (dist > light.distance) {
        return vec3(0);
    }

    vec3 L = normalize(lightPosition - viewPos);
    vec3 H = normalize(V + L);

    float attenuation = 1.0 / (light.constant + light.linear * dist + light.quadratic * (dist * dist));
    return calcCookTorrance(H, V, N, material, F0, L, light.color * attenuation);
}

vec3 CalcSpot(vec3 F0, vec3 lightPosition, Light light, Mtrl material, vec3 N, vec3 viewPos, vec3 V) {

    vec3 viewLightDirection = (view * vec4(light.direction, 0)).xyz;
    float dist = distance(lightPosition, viewPos);
    if (dist > light.distance) {
        return vec3(0);
    }

    vec3 L = normalize(lightPosition - viewPos);
    vec3 H = normalize(V + L);

    if (light.cutoff > 0) {
        float theta = dot(L, normalize(-viewLightDirection));
        if (theta < light.cutoff) { return vec3(0); }
    }
    float attenuation = 1.0 / (light.constant + light.linear * dist + light.quadratic * (dist * dist));

    return calcCookTorrance(H, V, N, material,F0, L, light.color * attenuation);
}

vec3 CalcDirectional(vec3 F0, vec3 lightPosition, Light light, Mtrl material, vec3 N, vec3 viewPos, vec3 V) {
    vec3 viewLightDirection = (view * vec4(light.direction, 0)).xyz;
    float dist = distance(lightPosition, viewPos);
    vec3 L = viewLightDirection;
    vec3 H = normalize(V + L);
    return calcCookTorrance(H, V, N, material, F0, L, light.color);
}
