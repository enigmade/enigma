import QtQuick
import QtQuick.Controls

// Runtimes tab (SPEC §18): one row per installed language runtime, with the
// active/default one flagged. Data comes from api.State.Runtimes.
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
                Label { text: modelData.name; font.bold: true; width: 120 }
                Label { text: modelData.version }
                Label {
                    text: modelData.active ? qsTr("(default)") : ""
                    color: "#3fb950"
                }
            }
        }
    }

    Label {
        anchors.centerIn: parent
        visible: parent.model.length === 0
        text: qsTr("No runtimes reported by daemon")
        opacity: 0.6
    }
}
