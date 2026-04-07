//+vert
#version 120

in vec3 a_Pos;
in vec3 a_Norm;
in vec2 a_TexCoord0;
// layout (location = 0) in vec3 a_Pos;
// layout (location = 1) in vec3 a_Norm;
// layout (location = 2) in vec2 a_TexCoord0;

out vec2 v_TexCoord0;

uniform mat4 u_VPMatrix;
uniform mat4 u_MMatrix;

void main() {
    gl_Position = u_VPMatrix * u_MMatrix * vec4(a_Pos.xyz, 1.0);
    v_TexCoord0 = a_TexCoord0;
}
//+frag
#version 120
precision mediump float;

out vec4 f_Color;

in vec2 v_TexCoord0;

uniform sampler2D u_DiffTexture;
/*
uniform vec3 u_DirLightDir[2];
uniform vec3 u_DirLightColor[2];
uniform float u_DirLightPower[2];

uniform vec3 u_PtLightPos[8];
uniform vec3 u_PtLightColor[8];
uniform float u_PtLightPower[8];
*/
void main() {
    f_Color = texture2D(u_DiffTexture, v_TexCoord0);
}