#include <QApplication>
#include <QLabel>
#include <QVBoxLayout>
#include <QWidget>

int main(int argc, char* argv[]) {
    QApplication app(argc, argv);

    QWidget window;
    window.setWindowTitle("Hello Qt6");
    window.setFixedSize(400, 200);

    auto* layout = new QVBoxLayout(&window);
    auto* label = new QLabel("Hello from Qt6!");
    label->setAlignment(Qt::AlignCenter);
    label->setStyleSheet("font-size: 24px; font-weight: bold;");
    layout->addWidget(label);

    window.show();
    return app.exec();
}
