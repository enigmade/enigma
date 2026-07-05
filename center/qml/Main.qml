import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

import "tabs"

// Enigma Center — the visual face of the enigma CLI (SPEC §18).
// Thin client: every tab renders the `daemon` context property (a
// DaemonClient bound to GET /v1/state). No business logic lives here.
ApplicationWindow {
    id: root
    width: 960
    height: 640
    visible: true
    title: qsTr("Enigma Center")

    // Poll the daemon so the UI stays live as services start/stop.
    Timer {
        interval: 3000
        running: true
        repeat: true
        onTriggered: daemon.refresh()
    }

    header: ToolBar {
        RowLayout {
            anchors.fill: parent
            Label {
                text: qsTr("Enigma Center")
                font.bold: true
                Layout.leftMargin: 12
            }
            Item { Layout.fillWidth: true }
            Label {
                text: daemon.online ? qsTr("● daemon online") : qsTr("○ daemon offline")
                color: daemon.online ? "#3fb950" : "#f85149"
                Layout.rightMargin: 12
            }
        }
    }

    ColumnLayout {
        anchors.fill: parent
        spacing: 0

        TabBar {
            id: tabs
            Layout.fillWidth: true
            TabButton { text: qsTr("Runtimes") }
            TabButton { text: qsTr("Services") }
            TabButton { text: qsTr("AI") }
            TabButton { text: qsTr("Projects") }
            TabButton { text: qsTr("Windows") }
            TabButton { text: qsTr("System") }
        }

        StackLayout {
            Layout.fillWidth: true
            Layout.fillHeight: true
            currentIndex: tabs.currentIndex

            RuntimesTab { model: daemon.runtimes }
            ServicesTab { model: daemon.services }
            AITab       { model: daemon.models }
            ProjectsTab { model: daemon.projects }
            WindowsTab  {}
            SystemTab   { gpu: daemon.gpu }
        }
    }

    footer: ToolBar {
        visible: daemon.lastError.length > 0
        Label {
            text: daemon.lastError
            color: "#f85149"
            leftPadding: 12
        }
    }
}
