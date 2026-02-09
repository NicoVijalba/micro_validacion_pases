CREATE TABLE IF NOT EXISTS records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    emision DATETIME NOT NULL,
    nave VARCHAR(150) NOT NULL,
    viaje VARCHAR(100) NOT NULL,
    cliente VARCHAR(200) NOT NULL,
    booking VARCHAR(100) NOT NULL,
    rama ENUM('internacional', 'nacional') NOT NULL,
    contenedor VARCHAR(120) NOT NULL,
    puerto_descargue VARCHAR(150) NOT NULL,
    libre_retencion_hasta DATE NOT NULL,
    dias_libre INT NOT NULL DEFAULT 0,
    transportista VARCHAR(200) NOT NULL DEFAULT '',
    titulo_terminal VARCHAR(200) NOT NULL,
    usuario_firma VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_booking_viaje_contenedor (booking, viaje, contenedor)
);
