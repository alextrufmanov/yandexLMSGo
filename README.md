Учебный проект YandexLMS "Программирование на Go | 24".
Спринт 1. Финальная задача 1.13.

ОБЩАЯ ИНФОРМАЦИЯ

Веб-сервис позволяет пользователю отправляет арифметическое выражение по HTTP (REST API) и получает в ответ результат его вычисления.
Поддерживает 4 основные арифметические операции + - * /. Выражение не должно содержать знак =.

Ожидает POST запрос на эндпоинт /api/v1/calculate
Body запроса должно содержать JSON следующего вида:
{
    "expression": "Выражение"  
}

Выражение может содержать следующие символы ()+-*/0123456789 и должно быть корректным математическим выражением.

В случае успешного вычисления выражения сервис возвращает код 200, в качестве результата в body возвращает JSON следующего вида:
{
    "result": "Результат"
}

В случае если входные данные не соответствуют требованиям, сервис возвращает код 422, а в body JSON следующего вида:
{
    "error": "Expression is not valid"
}

Если в работе сервиса произошла какая-то ошибка, код 500 и в body JSON следующего вида:
{
    "error": "Internal server error"
}

ЗАПУСК СЕРВИСА.

Сервис обрабатывает запросы, поступающие на порт 80.

Для его запуска в ОС Windows необходимо :
1. нажмите Win+R, либо перейдите в меню «Пуск» и выбрать пункт «Выполнить».
2. В появившемся окне набрать cmd и нажать ОК. 
3. В открывшемся окне (это и есть командная строка Windows) перейти в каталог проекта выполнив команду
   cd полныйПутьККакталогуПроекта
   например cd d:\Projects\yandexLMSGo. Если при этом надо будет сменить диск, предварительно введите название этого диска с двоеточием и нажмите Enter.
4. Для запуска сервиса введите в командной строке 
   go run main.go
   и нажмие Enter.


Проверка работы сервиса c помощью curl (примеры вызовов REST API с результатами):

1. Простое выражение 2+2*2. Наберите в командной строке
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"2+2*2\"}" localhost:80/api/v1/calculate
Результат:
{"result":"6"}

2. Простое выражение -2+2*2 (первый унарный минус). Наберите в командной строке
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"-2+2*2\"}" localhost:80/api/v1/calculate
Результат:
{"result":"2"}

3. Простое выражение 2+2*(-2) (унарный минус в скобках). Наберите в командной строке
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"2+2*(-2)\"}" localhost:80/api/v1/calculate
Результат:
{"result":"-2"}

4. Простое выражение со скобками (2+2)*2 (унарный минус в скобках). Наберите в командной строке
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"(2+2)*2\"}" localhost:80/api/v1/calculate
Результат:
{"result":"8"}

5. Сложное выражение -3*(12/4+(-2+8/2))+(7-20/5*(9-2*2*2+1))*(-10). Наберите в командной строке
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"-3*(12/4+(-2+8/2))+(7-20/5*(9-2*2*2+1))*(-10)\"}" localhost:80/api/v1/calculate
Результат:
{"result":"-5"}

6. Сложное выражение -(3*(12/4+(-2+8/2))+(7-20/5*(9-2*2*2+1))*(-10)). Наберите в командной строке
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"-(3*(12/4+(-2+8/2))+(7-20/5*(9-2*2*2+1))*(-10))\"}" localhost:80/api/v1/calculate
{"result":"-25"}

7. Сложное выражение с пробелами -(3*( 12/ 4+(-2 + 8/2 ) ) +(7-20/5*(9-2* 2*2+1))*(  -10)). Наберите в командной строке
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"-(3*( 12/ 4+(-2 + 8/2 ) ) +(7-20/5*(9-2* 2*2+1))*(  -10))\"}" localhost:80/api/v1/calculate
{"result":"-25"}

8. Ошибка в выражении (непарная скобка) -(3*(12/4+(-2+8/2))+(7-20/5*(9-2*2*2+1))*(-10). Наберите в командной строке
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"-3*(12/4+(-2+8/2))+(7-20/5*(9-2*2*2+1))*(-10)\"}" localhost:80/api/v1/calculate
Результат:
{
    "error": "Expression is not valid"
}

9. Ошибка в выражении (недопустимый символ =) -(3*(12/4+(-2+8/2))+(7-20/5*(9-2*2*2+1))*(-10))=. Наберите в командной строке
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"-3*(12/4+(-2+8/2))+(7-20/5*(9-2*2*2+1))*(-10))=\"}" localhost:80/api/v1/calculate
Результат:
{
    "error": "Expression is not valid"
}

10. Ошибка в выражении (недопустимый символы) -(3*(12/4+(-A+8/2))+(7-20/5*(9-2*2*2+b))*(-10))=. Наберите в командной строке
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"-3*(12/4+(-A+8/2))+(7-20/5*(9-2*2*2+b))*(-10))\"}" localhost:80/api/v1/calculate
Результат:
{
    "error": "Expression is not valid"
}

11. Ошибка в хосте. Наберите в командной строке
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"2+2*2\"}" localhostttt:80/api/v1/calculate
Результат:
curl: (6) Could not resolve host: localhostttt

12. Ошибка в номере порта. Наберите в командной строке
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"2+2*2\"}" localhost:801/api/v1/calculate
Результат:
curl: (7) Failed to connect to localhost port 801 after 2255 ms: Could not connect to server

13. Ошибка в пути. Наберите в командной строке
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"2+2*2\"}" localhost:80/api/v2/recalculate
Результат:
404 page not found

14. Ошибка в типе метода запроса. Наберите в командной строке
curl -X GET -H "Content-Type: application/json" -d "{\"expression\": \"2+2*2\"}" localhost:80/api/v1/calculate
Результат:
404 page not found

15. Ошибка в содержимом body (не хватает двойных первых кавычек) - JSON невозможно корректго распарсить. Наберите в командной строке
curl -X POST -H "Content-Type: application/json" -d "{expression\": \"2+2*2\"}" localhost:80/api/v1/calculate
Результат:
{
    "error": "Internal server error"
}
