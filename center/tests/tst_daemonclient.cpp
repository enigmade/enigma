#include "daemonclient.h"

#include <QSignalSpy>
#include <QTcpServer>
#include <QTcpSocket>
#include <QTest>
#include <QUrl>

// MockDaemon is a minimal HTTP/1.1 server that answers GET /v1/state with a
// canned api.State payload, standing in for the real enigma daemon so the
// client can be tested headlessly (SPEC §18).
class MockDaemon : public QObject
{
    Q_OBJECT
public:
    explicit MockDaemon(QByteArray body, QObject *parent = nullptr)
        : QObject(parent), m_body(std::move(body))
    {
        connect(&m_server, &QTcpServer::newConnection, this, &MockDaemon::onConnection);
        m_server.listen(QHostAddress::LocalHost, 0);
    }

    QUrl url() const
    {
        return QUrl(QStringLiteral("http://127.0.0.1:%1").arg(m_server.serverPort()));
    }

private slots:
    void onConnection()
    {
        QTcpSocket *sock = m_server.nextPendingConnection();
        connect(sock, &QTcpSocket::readyRead, this, [this, sock]() {
            sock->readAll(); // ignore the request line; always answer /v1/state
            QByteArray resp = "HTTP/1.1 200 OK\r\n";
            resp += "Content-Type: application/json\r\n";
            resp += "Content-Length: " + QByteArray::number(m_body.size()) + "\r\n";
            resp += "Connection: close\r\n\r\n";
            resp += m_body;
            sock->write(resp);
            sock->disconnectFromHost();
        });
    }

private:
    QTcpServer m_server;
    QByteArray m_body;
};

class TestDaemonClient : public QObject
{
    Q_OBJECT

private slots:
    void populatesStateFromDaemon();
    void reportsErrorWhenUnset();
    void handlesInvalidJson();
};

void TestDaemonClient::populatesStateFromDaemon()
{
    const QByteArray state = R"({
        "runtimes": [{"name": "go", "version": "1.26", "active": true}],
        "services": [{"name": "enigma-ollama", "status": "running", "port": 11434}],
        "projects": [{"path": "/home/u/app", "port": 8080, "url": "https://app.test"}],
        "models": [{"name": "llama3.1:8b", "size_gb": 4.7, "backend": "ollama"}],
        "gpu": {"vendor": "NVIDIA", "model": "RTX 4090", "vram_gib": 24}
    })";

    MockDaemon daemon(state);
    DaemonClient client;
    client.setBaseUrl(daemon.url());

    QSignalSpy spy(&client, &DaemonClient::stateChanged);
    client.refresh();
    QVERIFY(spy.wait(2000));

    QCOMPARE(client.runtimes().size(), 1);
    QCOMPARE(client.runtimes().first().toMap().value("name").toString(), QStringLiteral("go"));
    QCOMPARE(client.services().size(), 1);
    QCOMPARE(client.projects().size(), 1);
    QCOMPARE(client.models().size(), 1);
    QCOMPARE(client.gpu().value("vram_gib").toInt(), 24);
    QVERIFY(client.online());
}

void TestDaemonClient::reportsErrorWhenUnset()
{
    DaemonClient client; // no base URL
    QSignalSpy errSpy(&client, &DaemonClient::errorChanged);
    client.refresh();
    QVERIFY(errSpy.count() >= 1);
    QVERIFY(!client.online());
    QVERIFY(!client.lastError().isEmpty());
}

void TestDaemonClient::handlesInvalidJson()
{
    MockDaemon daemon(QByteArray("this is not json"));
    DaemonClient client;
    client.setBaseUrl(daemon.url());

    QSignalSpy errSpy(&client, &DaemonClient::errorChanged);
    client.refresh();
    QVERIFY(errSpy.wait(2000));
    QVERIFY(client.lastError().contains("invalid state JSON"));
}

QTEST_MAIN(TestDaemonClient)
#include "tst_daemonclient.moc"
