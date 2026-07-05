#include "daemonclient.h"

#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QUrl>

#include <cstdlib>

// Resolves the daemon socket base URL. In production the enigma daemon
// listens on an HTTP-over-unix-socket at $XDG_RUNTIME_DIR/enigma-center.sock;
// ENIGMA_CENTER_URL overrides it (used to point at a mock in dev/CI).
static QUrl daemonBaseUrl()
{
    if (const char *override = std::getenv("ENIGMA_CENTER_URL")) {
        return QUrl(QString::fromUtf8(override));
    }
    const char *runtimeDir = std::getenv("XDG_RUNTIME_DIR");
    const QString sock = QString::fromUtf8(runtimeDir ? runtimeDir : "/tmp")
                         + QStringLiteral("/enigma-center.sock");
    // Qt addresses unix-socket HTTP via the unix+http scheme.
    return QUrl(QStringLiteral("unix+http://") + sock);
}

int main(int argc, char *argv[])
{
    QGuiApplication app(argc, argv);
    app.setApplicationName(QStringLiteral("Enigma Center"));

    DaemonClient client;
    client.setBaseUrl(daemonBaseUrl());
    client.refresh();

    QQmlApplicationEngine engine;
    engine.rootContext()->setContextProperty(QStringLiteral("daemon"), &client);
    engine.loadFromModule("EnigmaCenter", "Main");

    if (engine.rootObjects().isEmpty()) {
        return -1;
    }

    return app.exec();
}
