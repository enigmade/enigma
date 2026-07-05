#pragma once

#include <QObject>
#include <QUrl>
#include <QVariantMap>

class QNetworkAccessManager;

// DaemonClient is the ONLY bridge between the QML UI and the enigma daemon.
// It performs GET /v1/state (over the Unix socket in production, or a plain
// TCP URL in tests) and exposes the parsed api.State snapshot as QML
// properties. The UI holds no business logic — it renders whatever this
// object publishes (SPEC §18 thin-client rule).
class DaemonClient : public QObject
{
    Q_OBJECT
    Q_PROPERTY(QVariantList runtimes READ runtimes NOTIFY stateChanged)
    Q_PROPERTY(QVariantList services READ services NOTIFY stateChanged)
    Q_PROPERTY(QVariantList projects READ projects NOTIFY stateChanged)
    Q_PROPERTY(QVariantList models READ models NOTIFY stateChanged)
    Q_PROPERTY(QVariantMap gpu READ gpu NOTIFY stateChanged)
    Q_PROPERTY(bool online READ online NOTIFY onlineChanged)
    Q_PROPERTY(QString lastError READ lastError NOTIFY errorChanged)

public:
    explicit DaemonClient(QObject *parent = nullptr);

    // baseUrl points at the daemon. In production this is a unix+http URL;
    // tests pass an http://127.0.0.1:<port> mock. Settable so the same
    // binary can be pointed at a mock in headless CI.
    Q_INVOKABLE void setBaseUrl(const QUrl &url);

    QVariantList runtimes() const { return m_runtimes; }
    QVariantList services() const { return m_services; }
    QVariantList projects() const { return m_projects; }
    QVariantList models() const { return m_models; }
    QVariantMap gpu() const { return m_gpu; }
    bool online() const { return m_online; }
    QString lastError() const { return m_lastError; }

public slots:
    // Fetches GET /v1/state and updates the exposed properties.
    void refresh();

signals:
    void stateChanged();
    void onlineChanged();
    void errorChanged();

private:
    void applyState(const QByteArray &json);
    void setOnline(bool online);
    void setError(const QString &err);

    QNetworkAccessManager *m_net;
    QUrl m_baseUrl;
    QVariantList m_runtimes;
    QVariantList m_services;
    QVariantList m_projects;
    QVariantList m_models;
    QVariantMap m_gpu;
    bool m_online = false;
    QString m_lastError;
};
