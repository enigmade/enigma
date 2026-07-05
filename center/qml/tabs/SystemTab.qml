import QtQuick
import QtQuick.Controls

// System tab (SPEC §18): snapshots + rollback, staged updates, boot-time
// breakdown, doctor results. GPU comes from api.State.GPU (hardware.toml).
Item {
    property var gpu: ({})

    Column {
        anchors.fill: parent
        anchors.margins: 16
        spacing: 12

        Label { text: qsTr("System"); font.bold: true; font.pixelSize: 18 }

        GroupBox {
            title: qsTr("Detected GPU")
            width: parent.width
            Column {
                spacing: 4
                Label {
                    text: gpuField("vendor", qsTr("unknown")) + " " + gpuField("model", "")
                    font.bold: true
                }
                Label {
                    text: qsTr("VRAM: %1 GiB").arg(gpuField("vram_gib", 0))
                    opacity: 0.7
                }
            }
        }

        GroupBox {
            title: qsTr("Maintenance")
            width: parent.width
            Row {
                spacing: 8
                Button { text: qsTr("Rollback to snapshot") }
                Button { text: qsTr("Install staged update (~20s)") }
                Button { text: qsTr("Run doctor") }
            }
        }
    }

    function gpuField(key, fallback) {
        if (gpu && gpu[key] !== undefined)
            return gpu[key];
        return fallback;
    }
}
