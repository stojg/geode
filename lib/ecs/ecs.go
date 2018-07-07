package ecs

// Components are only raw data, ie component Position.x, Position.y ie a struct nothing else/more
// Entities is a collection of Components, Position, Motion, Input nothing else/more
// System, takes a list of Components, ie. All the Position and Motion component in the world, nothing more/else

// update systems
// for each system
// find out what allComponents it needs
// if a single component
//    just grab that component list and pass to it
// else
//    for all nextEntityID
//       if entity doesn't have all component
// 			skip
//       else
//          add allComponents to list
//    end
//      call system with component list
//
