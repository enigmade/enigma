import QtQuick
import QtQuick.Controls

// Projects tab (SPEC §18): detected dev projects with their .test URL and
// allocated port. Data comes from api.State.Projects.
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
            contentItem: Column {
                spacing: 2
                Label { text: modelData.path; font.bold: true }
                Row {
                    spacing: 12
                    Label { text: qsTr("port %1").arg(modelData.port); opacity: 0.7 }
                    Label {
                        text: modelData.url
                        color: "#58a6ff"
                        visible: modelData.url !== undefined && modelData.url.length > 0
                    }
                }
            }
        }
    }

    Label {
        anchors.centerIn: parent
        visible: parent.model.length === 0
        text: qsTr("No active projects")
        opacity: 0.6
    }
}
