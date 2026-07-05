import QtQuick
import QtQuick.Controls

// Windows apps tab (SPEC §18): Bottles + VM state (Tier1/Tier2), launch/add.
// The daemon's Windows endpoint is not yet in api.State, so this panel
// currently documents the intended layout and the DirectX/VC++ redistributable
// status surfaced by `enigma win setup`.
Item {
    Column {
        anchors.centerIn: parent
        spacing: 8
        Label {
            text: qsTr("Windows Apps (Wine Tier 1)")
            font.bold: true
            font.pixelSize: 18
            anchors.horizontalCenter: parent.horizontalCenter
        }
        Label {
            text: qsTr("Managed via Bottles. Required DirectX/VC++ redistributables:")
            opacity: 0.7
            anchors.horizontalCenter: parent.horizontalCenter
        }
        Label {
            text: qsTr("d3dcompiler_47 · d3dx9 · d3dx11_43 · vcrun2019")
            font.family: "monospace"
            anchors.horizontalCenter: parent.horizontalCenter
        }
    }
}
