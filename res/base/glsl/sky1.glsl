//+vert
#version 330 core

layout (location = 0) in vec2 a_Pos;

out vec4 v_Pos;

void main() {
    gl_Position = vec4(a_Pos.xy, 1.0, 1.0);
    v_Pos = vec4(a_Pos.xy, 0.0, 1.0);
}
//+frag
#version 330 core
precision mediump float;

out vec4 f_Color;

in vec4 v_Pos;

uniform mat4 u_PMatrix;
uniform mat4 u_VMatrix;

uniform vec3 u_SunDir;
uniform vec3 u_SunColor;

const vec3 sunInnerColor = vec3(1.00, 1.00, 1.00);
const vec3 sunOuterColor = vec3(0.98, 1.00, 0.65);

const vec3 skyDeepBlueColor = vec3(0.0, 0.4, 1.0);
const vec3 skyHazeBlueColor = vec3(0.0, 0.8, 1.0);
const vec3 skyDuskColor = vec3(1.0, 0.3, 0.1);

const float sunDownOffset = 6.0;
const float sunDuskInvar = 3.0;

vec3 duskDown(vec3 dirSun) {
    vec3 down = dirSun;
    down.z -= sunDownOffset;
    return -normalize(down);
}

vec3 sunColor(float isSun, vec3 dirSun) {
    vec3 color = mix(sunInnerColor, sunOuterColor, isSun);
    return color;
}

vec3 skyColor(float z, float zSun, float distDusk) {
    vec3 color = mix(skyHazeBlueColor, skyDeepBlueColor, sqrt(z));
    float duskInt = exp(-sunDuskInvar * zSun * zSun);
    color = mix(color, skyDuskColor, distDusk * duskInt);
    return color;
}

void main() {
    mat4 vd = u_VMatrix;
    vd[3] = vec4(0.0, 0.0, 0.0, 1.0);
    vec4 back = inverse(u_PMatrix * vd) * v_Pos;
    vec3 view = normalize(back.xyz / back.w);
    float z = max(view.z, 0.0);
    float zSun = u_SunDir.z;
    float distSun = distance(view, normalize(u_SunDir));
    float isSun = max(distSun * 100.0 - 0.5, 0.0);
    float distDusk = distance(view, duskDown(u_SunDir));
    float isDusk = max(distDusk - 0.65, 0.0);
    float sunMask = clamp(isSun * isSun, 0.0, 1.0);
    float sunInt = length(u_SunColor) / 1.73;
    vec3 color = sunInt * skyColor(z, zSun, isDusk);
    color = mix(sunColor(isSun, u_SunDir), color, sunMask);
    f_Color = vec4(color, 1.0);
}