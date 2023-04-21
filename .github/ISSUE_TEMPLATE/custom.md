---
name: Custom issue template
about: Describe this issue template's purpose here.
title: ''
labels: ''
assignees: ''

---

Se requiere realizar mantenimiento a la funcionalidad de notificaciones a nivel de sistema y por correo electrónico.

Actualmente se estan presentando los siguientes errores;

- No llega el link de suscripción por correo electrónico.
- No llega confirmación de suscripción por correo electrónico.
- La validación de suscripción no funciona correctamente ya que pide a usuarios ya suscritos que vuelvan a suscribirse.

**Especificaciones técnicas**

- Las especificaciones técnicas del modulo de notificaciones se encuentra en su [documento técnico](https://drive.google.com/file/d/1ZlXdKAyooyftamyCZX8Btmd5nhDR5YEX/view?usp=sharing) relacionado.

**Sub Tareas**

- [ ] Revisión de [documentación existente](https://drive.google.com/file/d/1ZlXdKAyooyftamyCZX8Btmd5nhDR5YEX/view?usp=sharing).
- [ ] Revisión de[ pruebas funcionales](https://docs.google.com/spreadsheets/d/1SuAhfDHYM0OxwZ1mrV-MM36aEN5v8Cia/edit?usp=sharing&ouid=104124925692794756698&rtpof=true&sd=true) realizadas.
- [ ] Acceder a AWS y revisar funcionamiento del servicio de notificaciones.
- [ ] Instalación pruebas del API notificaciones_mid en ambiente local.
- [ ] Revisión y análisis de posibles problemas.
- [ ] Solución de problemas (desarrollo).

**Criterios de aceptación**

- [ ] A cada usuario que no este suscrito a las notificaciones del sistema, le debe llegar un mensaje de suscripción al correo que tiene inscrito en WSo2.
- [ ] Cuando el usuario se suscriba mediante el link, su correo debe quedar asociado al topic que maneja el sistema, de tal forma que cuando se ejecute una acción de notificación asociada al topic, se le envie la notificación por correo.
- [ ] La suscripción de un usuario debe mantenerse en el tiempo y no caducar, hasta que se indique lo contrario.
- [ ] Verificar que las notificaciones internas del sistema esten funcionando correctamente.

**Requerimientos**

- [ ] Listado de variables de entorno de notificaciones_mid.
- [ ] Usuario con acceso a servicio de notificaciones en AWS.

**Definition of Ready - DoR**

- [x] Está refinada y estimada en puntos de historia por el equipo.
- [x] Incluye la descripción y criterios de aceptación, con el detalle funcional y especificaciones técnicas, de forma entendible por cualquier miembro del equipo.
- [ ] No tiene bloqueos que impidan su ejecución.
- [ ] Las dependencias entán identificadas y resueltas.
- [ ] Puede ser probada dentro del Sprint. 

**Definition of Done - DoD - Desarrollo**

- [ ] Desarrollo en local.
- [ ] Push en Feature.
- [ ] Pruebas locales (funcionales).
- [ ] PR a Develop.
- [ ] Criterios de aceptación cumplidos.
- [ ] Documentación de issue realizada.
- [ ] Aprobada por SM/Líder técnico.
