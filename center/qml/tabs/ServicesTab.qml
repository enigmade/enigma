import QtQuick
import QtQuick.Controls

// Services tab (SPEC §18): every enigma-managed service with status, port,
// and a start/stop toggle. Data comes from api.State.Services. The toggle
// posts an action back to the daemon (wired via daemon.startService/stop
// once those endpoints land).
Item {
    property var model: []

    ListView {
        anchors.fill: parent
        anchors.margins: 12
        spacing: 6
        model: parent.model

        delegate: ItemDelegate {
            required property var modelData
            width: ListView.view.width
            contentItem: Row {
                spacing: 12
                Rectangle {
                    width: 10; height: 10; radius: 5
                    anchors.verticalCenter: parent.verticalCenter
                    color: modelData.status === "running" ? "#3fb950" : "#8b949e"
                }
                Label { text: modelData.name; font.bold: true; width: 200 }
                Label { text: modelData.status; width: 90 }
                Label {
                    text: modelData.port > 0 ? (":" + modelData.port) : ""
                    opacity: 0.7
                }
            }
        }
    }

    Label {
        anchors.centerIn: parent
        visible: parent.model.length === 0
        text: qsTr("No services running")
        opacity: 0.6
    }
}
