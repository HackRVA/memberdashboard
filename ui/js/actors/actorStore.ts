interface Actors {
  [address: string]: any;
}
class ActorStore {
  actors: Actors = {};

  /**
   * register an actor and return it's address
   * @param name
   * @param actorType
   */
  register(name: string, actorType: any): String {
    this.actors[name] = new actorType();
    console.log("successfully registered actor: ", name);
    return name;
  }

  /**
   * lookup an actor to pass it messages
   * @param address
   */
  lookup(address: string): any {
    return this.actors[address];
  }
}

export default new ActorStore();
