import { Notice, Plugin } from 'obsidian';
import { getService } from '../../sbe-core/src/bridge';
import { errorMessage } from '../../sbe-core/src/utils/errors';

export default class ExamplePlugin extends Plugin {
  async onload(): Promise<void> {
    this.addCommand({
      id: 'example-check-llm',
      name: 'Проверить LLM через стенд',
      callback: async () => {
        try {
          const dev = await getService('sbe-devstand');
          const me = await dev.whoami();
          const answer = await dev.llm.complete(
            'Ты кратко отвечаешь по-русски.',
            'Назови три свойства фасадной краски.',
          );
          new Notice(`Стенд: ${me.email}. Ответ LLM: ${answer.slice(0, 120)}`);
        } catch (e) {
          new Notice(errorMessage(e));
        }
      },
    });
  }
}
