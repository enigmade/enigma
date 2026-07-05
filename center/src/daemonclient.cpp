#include "daemonclient.h"

#include <QJsonArray>
#include <QJsonDocument>
#include <QJsonObject>
#include <QNetworkAccessManager>
#include <QNetworkReply>
#include <QNetworkRequest>

DaemonClient::DaemonClient(QObject *parent)
    : QObject(parent)
    , m_net(new QNetworkAccessManager(this))
{
}

void DaemonClient::setBaseUrl(const QUrl &url)
{
    m_baseUrl = url;
}

void DaemonClient::refresh()
{
    if (m_baseUrl.isEmpty()) {
        setError(QStringLiteral("daemon base URL not set"));
        setOnline(false);
        return;
    }

    QUrl stateUrl = m_baseUrl;
    stateUrl.setPath(QStringLiteral("/v1/state"));

    QNetworkRequest req(stateUrl);
    QNetworkReply *reply = m_net->get(req);

    connect(reply, &QNetworkReply::finished, this, [this, reply]() {
        reply->deleteLater();
        if (reply->error() != QNetworkReply::NoError) {
            setError(reply->errorString());
            setOnline(false);
            return;
        }
        // A malformed body means the daemon answered but we can't use it:
        // keep the parse error set by applyState() and don't mark online.
        if (applyState(reply->readAll())) {
            setOnline(true);
            setError(QString());
        } else {
            setOnline(false);
        }
    });
}

bool DaemonClient::applyState(const QByteArray &json)
{
    QJsonParseError perr;
    const QJsonDocument doc = QJsonDocument::fromJson(json, &perr);
    if (perr.error != QJsonParseError::NoError || !doc.isObject()) {
        setError(QStringLiteral("invalid state JSON: %1").arg(perr.errorString()));
        return false;
    }

    const QJsonObject obj = doc.object();
    m_runtimes = obj.value(QStringLiteral("runtimes")).toArray().toVariantList();
    m_services = obj.value(QStringLiteral("services")).toArray().toVariantList();
    m_projects = obj.value(QStringLiteral("projects")).toArray().toVariantList();
    m_models = obj.value(QStringLiteral("models")).toArray().toVariantList();
    m_gpu = obj.value(QStringLiteral("gpu")).toObject().toVariantMap();

    emit stateChanged();
    return true;
}

void DaemonClient::setOnline(bool online)
{
    if (m_online != online) {
        m_online = online;
        emit onlineChanged();
    }
}

void DaemonClient::setError(const QString &err)
{
    if (m_lastError != err) {
        m_lastError = err;
        emit errorChanged();
    }
}
