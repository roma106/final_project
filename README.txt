
Интрефейс программы

            1. Окно ввода выражения
            2. База выражений
            3. Настройки операций
        
Как пользоваться программой
0. База данных
        Разверните базу данных с именем go_projects
        Создайте таблицу yandex_final со столбцами Expression(varchar), Status(varchar), Result(integer), StartingTime(timestamp with tz), EndingTime(timestamp with tz), ID (bigint primary key auto_increment). (sql скрипт прилагается снизу)
1. Запуск сервера
            Откройте терминал.  Перейдите в директорию backend (cd backend). Запустите сервер приложения командой go run main.go.
2. Веб-интерфейс программы
            Зайдите в папку frontend.  Откройте файл index.html.
3. Работа с программой
            Отконфигурируйте время выполнения операций в окне "настройка операций". Введите нужное выражение в поле ввода, оно должно содержать больше двух символов. Нажмите на кнопку справа (калькулятор). Окно "база выражений" обновлено, ожидайте вычисления вашего примера.
        


ТГ для связи: @Romanovski228

SQL скрипт для создания БД и таблицы:

-- Database: go_projects

-- DROP DATABASE IF EXISTS go_projects;

CREATE DATABASE go_projects
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'English_United States.1252'
    LC_CTYPE = 'English_United States.1252'
    LOCALE_PROVIDER = 'libc'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;



-- Table: public.yandex_final

-- DROP TABLE IF EXISTS public.yandex_final;

CREATE TABLE IF NOT EXISTS public.yandex_final
(
    "Expression" character varying COLLATE pg_catalog."default" NOT NULL,
    "Status" character varying COLLATE pg_catalog."default" NOT NULL,
    "Result" integer,
    "StartingTime" timestamp with time zone NOT NULL,
    "EndingTime" timestamp with time zone NOT NULL,
    "ID" bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 1000 CACHE 1 ),
    CONSTRAINT yandex_final_pkey PRIMARY KEY ("ID")
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.yandex_final
    OWNER to postgres;
