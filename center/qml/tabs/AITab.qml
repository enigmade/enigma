import QtQuick
import QtQuick.Controls

// AI tab (SPEC §18): installed models with sizes and a combined VRAM meter.
// Data comes from api.State.Models. "Won't fit" warnings (§20.10) compare
// total model size against detected VRAM before launch.
Item {
    property var model: []

    Column {
        anchors.fill: parent
        anchors.margins: 12
        spacing: 12

        Label {
            text: qsTr("Total on disk: %1 GB").arg(totalSizeGB().toFixed(1))
            font.bold: true
        }

        ListView {
            width: parent.width
            height: parent.height - 40
            spacing: 6
            model: parent.parent.model

            delegate: ItemDelegate {
                required property var modelData
                width: ListView.view.width
                contentItem: Row {
                    spacing: 12
                    Label { text: modelData.name; font.bold: true; width: 240 }
                    Label { text: (modelData.size_gb).toFixed(1) + " GB"; width: 90 }
                    Label { text: modelData.backend; opacity: 0.7 }
                }
            }
        }
    }

    Label {
        anchors.centerIn: parent
        visible: parent.model.length === 0
        text: qsTr("No AI models installed")
        opacity: 0.6
    }

    function totalSizeGB() {
        var total = 0;
        for (var i = 0; i < model.length; ++i)
            total += model[i].size_gb;
        return total;
    }
}
