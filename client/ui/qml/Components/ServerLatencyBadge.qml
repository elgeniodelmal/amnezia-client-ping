import QtQuick
import QtQuick.Layouts

import ServerQuality 1.0
import "../Controls2"
import "../Controls2/TextTypes"

RowLayout {
    id: root

    property int latencyMs: -1
    property int quality: ServerQuality.Unknown
    property bool testing: false

    visible: testing || quality !== ServerQuality.Unknown

    SpinnerType {
        visible: root.testing
        running: visible
        size: SpinnerType.Small
    }

    CaptionTextType {
        visible: !root.testing
        text: root.latencyMs >= 0 ? qsTr("%1 ms").arg(root.latencyMs) : qsTr("No reply")
        wrapMode: Text.NoWrap
    }
}
