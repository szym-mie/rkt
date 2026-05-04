//+vert
#version 330 core

layout (location = 0) in vec3 a_Pos;
out vec4 v_Pos;

void main() {
    v_Pos = vec4(a_Pos.xyz, 1.0);
    gl_Position = vec4(a_Pos.xy, 1.0, 1.0);
}
//+frag
precision mediump float;

out vec4 f_Color;

in vec4 v_Pos;

uniform mat4 u_VPIMatrix;
uniform samplerCube u_DiffTexture;

void main() {
    vec4 backVec = u_VPIMatrix * v_Pos;
    f_Color = texture(u_DiffTexture, backVec.xyz / backVec.w);
}