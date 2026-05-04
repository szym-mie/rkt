//+vert
#version 330 core

layout (location = 0) in vec3 a_Pos;
layout (location = 1) in vec3 a_Norm;
layout (location = 2) in vec3 a_Tang;
layout (location = 4) in vec2 a_UV0;

out vec3 v_FragPos;
out vec3 v_Norm;
out vec2 v_UV0;
out mat3 v_TBNMatrix;

uniform mat4 u_PMatrix;
uniform mat4 u_VMatrix;
uniform mat4 u_MMatrix;

void main() {
    gl_Position = u_PMatrix * u_VMatrix * u_MMatrix * vec4(a_Pos.xyz, 1.0);
    mat3 mmat = mat3(u_MMatrix);
    v_Norm = mmat * normalize(a_Norm);
    v_FragPos = vec3(u_MMatrix * vec4(a_Pos, 1.0));
    v_UV0 = a_UV0;

    vec3 n = normalize(mmat * a_Norm);
    vec3 t = normalize(mmat * a_Tang);
    vec3 b = cross(n, t);
    v_TBNMatrix = mat3(t, b, n);
}
//+frag
#version 330 core
precision mediump float;

#define UV1_SCALE 16.0

out vec4 f_Color;

in vec3 v_FragPos;
in vec3 v_Norm;
in vec2 v_UV0;
in mat3 v_TBNMatrix;

uniform sampler2D u_DiffTexture0;
uniform sampler2D u_DiffTexture1;
uniform sampler2D u_NormTexture;
uniform vec3 u_AmbLightColor;
uniform vec3 u_DirLightDir[2];
uniform vec3 u_DirLightColor[2];
/*
uniform vec3 u_PtLightPos[8];
uniform vec3 u_PtLightColor[8];
uniform float u_PtLightPower[8];
*/
vec3 calcDirLight(vec3 dir, vec3 color, vec3 norm, vec3 tex0, vec3 tex1) {
    float diff = max(dot(dir, norm), 0.0);
    vec3 diffuse = color * diff * tex0 * tex1;
    return diffuse;
}

vec3 calcPtLight(vec3 pos, vec3 color, float linFall, float quaFall, vec3 norm, vec3 tex0, vec3 tex1) {
    vec3 dir = normalize(pos - v_FragPos);
    float diff = max(dot(dir, norm), 0.0);
    float dist = length(pos - v_FragPos);
    float atten = 1.0 / (1.0 + linFall * dist + quaFall * (dist * dist));
    vec3 diffuse = color * diff * tex0 * tex1;
    diffuse *= atten;
    return diffuse;
}

void main() {
    vec3 tex0 = texture(u_DiffTexture0, v_UV0).rgb;
    vec3 tex1 = texture(u_DiffTexture1, v_UV0 * UV1_SCALE).rgb;
    vec3 nmap = texture(u_NormTexture, v_UV0 * UV1_SCALE).rgb;
    nmap = normalize(nmap * 2.0 - 1.0);
    nmap = normalize(v_TBNMatrix * nmap);

    vec3 color = u_AmbLightColor * tex0 * tex1;
    for (int i = 0; i < 2; i++) {
        color += calcDirLight(u_DirLightDir[i], u_DirLightColor[i], nmap, tex0, tex1);
    }
    f_Color = vec4(color, 1.0);
}