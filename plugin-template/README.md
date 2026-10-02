# Шаблон плагина

Минимальный каркас плагина Obsidian для разработки под SBE.

Положите папку в `<вольт>/.obsidian/plugins/example-plugin/`, а рядом должна
лежать папка `sbe-core` (ядро системы) — сборка использует её относительными
импортами.

```bash
npm install
npm run build
```

Затем включите плагин в Obsidian («Сторонние плагины»).

В `src/main.ts` показано, как получить доступ к стенду:

```ts
const dev = await getService('sbe-devstand');
await dev.llm.complete(system, user);   // боевой LLM, только чтение/генерация
await dev.kb.search('запрос');          // боевая База знаний
const local = dev.local;                // { url, token } — локальная запись
```

Полное руководство по созданию плагина — в репозитории вольта,
`docs/sbe-new-plugin-guide.md`. Правила кода: без `any`, `catch(e: unknown)` +
`errorMessage()`, `requestUrl()`, `window.setTimeout()`, классы `tn-*`.
