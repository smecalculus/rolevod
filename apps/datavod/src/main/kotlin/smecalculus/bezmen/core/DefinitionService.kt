package smecalculus.bezmen.core

import java.util.UUID
import java.util.UUID.randomUUID
import kotlin.reflect.KClass
import kotlin.reflect.KProperty1
import kotlin.reflect.full.companionObjectInstance
import kotlin.reflect.full.declaredMemberFunctions
import kotlin.reflect.full.declaredMemberProperties
import kotlin.reflect.full.findAnnotation
import kotlin.reflect.full.isSupertypeOf
import kotlin.reflect.full.starProjectedType

class SessionDao {
    private val sessions: MutableMap<UUID, Session> = mutableMapOf()

    fun save(session: Session) {
        sessions[session.id] = session.copy()
    }

    fun update(session: Session) {
        val s = sessions[session.id]!!
        s.choices = session.choices
        s.message = session.message
    }

    fun get(id: UUID): Session {
        return sessions[id]!!
    }
}

class DefinitionService {

    private val types: MutableMap<String, KClass<*>> = mutableMapOf()
    private val sessionDao = SessionDao()

    fun register(typeClass: KClass<*>) {
        val annotation = typeClass.findAnnotation<Agent>()
        requireNotNull(annotation) {
            "Class '${typeClass.simpleName}' must have '${Agent::class.simpleName}' annotation"
        }
        types[annotation.name] = typeClass
    }

    fun create(name: String): Instance {
        val typeClass = types[name]
        requireNotNull(typeClass) {
            "There are no class with '${name}' name"
        }
        return Instance(randomUUID(), typeClass)
    }

    fun assign(whole: Instance, placement: String, part: Instance): Assignment {
        val kPlacement = whole.type.declaredMemberProperties
            .find { it.findAnnotation<Placement>()?.name == placement }
        requireNotNull(kPlacement) {
            "Class '${whole.type.simpleName}' has no '${placement}' placement"
        }
        require(kPlacement.returnType.isSupertypeOf(part.type.starProjectedType)) {
            "Class '${part.type.simpleName}' must implement '${kPlacement.returnType}' interface"
        }
        return Assignment(randomUUID(), whole, kPlacement, part)
    }

    fun initiate(assignment: Assignment, port: String): Session {
        val kRole = assignment.placement.returnType.classifier as KClass<*>
        val kPort = kRole.declaredMemberFunctions.find { it.findAnnotation<Port>()?.name == port }
        requireNotNull(kPort) {
            "Role '${kRole.simpleName}' has no '${port}' port"
        }
        require(kPort.returnType.classifier == Interplay::class) {
            "Method '${kPort.name}' must return '${Interplay::class.qualifiedName}'"
        }
        val interplay = kPort.call(kRole.companionObjectInstance) as Interplay
        val session = Session(randomUUID(), assignment.part.id, assignment.whole.id, mutableListOf(), Unit, interplay.tree.node)
        sessionDao.save(session)
        return session
    }

    fun <T : Any> send(sender: Instance, session: Session, message: T) {
        when (val state = session.state) {
            is TheirChoice<*> -> {
                require(sender.id == session.they) {
                    "Only '${session.they}' can send"
                }
                require(state.branches.containsKey(message::class)) {
                    "Allowed message classes: ${state.branches.keys}"
                }
                val selectedState = state.branches[message::class]!!
                session.state = selectedState
                session.choices.add(selectedState.name)
                session.message = message
                sessionDao.update(session)
            }

            is OurChoice<*> -> {
                require(sender.id == session.we) {
                    "Only '${session.we}' can send"
                }
                require(state.branches.containsKey(message::class)) {
                    "Allowed message classes: ${state.branches.keys}"
                }
                val selectedState = state.branches[message::class]!!
                session.state = selectedState
                session.choices.add(selectedState.name)
                session.message = message
                sessionDao.update(session)
            }

            else -> throw IllegalStateException("Send impossible in '${state::class.simpleName}' state")
        }
    }

    fun <T : Any> receive(receiver: Instance, session: Session, messageClass: KClass<T>): T {
        val message = session.message
        require(messageClass == message::class) {
            "Must receive '${message::class.qualifiedName}'"
        }
        return when (val state = session.state) {
            is WeFromThem<*> -> {
                require(receiver.id == session.we) {
                    "Only '${session.we}' can receive"
                }
                session.state = state.node
                session.message = Unit
                sessionDao.update(session)
                message as T
            }

            is TheyFromUs<*> -> {
                require(receiver.id == session.they) {
                    "Only '${session.they}' can receive"
                }
                session.state = state.node
                session.message = Unit
                sessionDao.update(session)
                message as T
            }

            else -> throw IllegalStateException("Receive impossible in '${state::class.simpleName}' state")
        }
    }

    fun restore(id: UUID): Session {
        val session = sessionDao.get(id)
        val choices = session.choices.toMutableList()
        while (choices.isNotEmpty()) {
            when (val state = session.state) {
                is OurChoice<*> -> {
                    val choice = choices.removeAt(0)
                    val messageClass = state.branches.keys.find { it.findAnnotation<State>()?.name == choice }
                    session.state = state.branches[messageClass]!!
                }
                is TheirChoice<*> -> {
                    val choice = choices.removeAt(0)
                    val messageClass = state.branches.keys.find { it.findAnnotation<State>()?.name == choice }
                    session.state = state.branches[messageClass]!!
                }
                else -> session.state = state.node
            }
        }
        return session;
    }
}

annotation class Agent(val name: String)
typealias Construction = Agent
typealias Module = Agent
typealias OrgUnit = Agent

annotation class Role(val name: String)
typealias Func = Role

annotation class Port(val name: String)
typealias Endpoint = Port

annotation class Placement(val name: String)

annotation class Service(val name: String)
annotation class Interface(val spec: KClass<*>)

annotation class Alpha(val name: String)
annotation class State(val name: String)
annotation class Artifact(val name: String)
typealias WorkProduct = Artifact

data class Instance(val id: UUID, val type: KClass<*>)

data class Assignment(val id: UUID, val whole: Instance, val placement: KProperty1<*, *>, val part: Instance)

data class Session(
    val id: UUID,
    val we: UUID,
    val they: UUID,
    var choices: MutableList<String>,
    var message: Any,
    var state: TreeNode
)
typealias Process = Session
typealias Project = Session
typealias Issue = Session
typealias Case = Session

data class Interplay(val name: String, val tree: TreeNode)
typealias Practice = Interplay
typealias Scenario = Interplay
typealias Protocol = Interplay

fun interplay(name: String, configure: Tree.() -> Unit): Interplay {
    val tree = Tree()
    tree.configure()
    return Interplay(name, tree)
}

@TreeMarker
sealed class TreeNode {
    lateinit var node: TreeNode

    inline fun <reified T : Any> their(configure: TheirChoice<T>.() -> Unit) {
        val tagClass = T::class
        val annotation = tagClass.findAnnotation<Alpha>()
        requireNotNull(annotation) {
            "Class '${tagClass.simpleName}' must have '${Alpha::class.simpleName}' annotation"
        }
        require(tagClass.isSealed) {
            "Class '${tagClass.simpleName}' must be sealed"
        }
        val choice = TheirChoice<T>()
        choice.configure()
        val uncovered = tagClass.sealedSubclasses.toSet() subtract choice.branches.keys
        require(uncovered.isEmpty()) {
            "There are uncovered classes: ${uncovered.map { it::simpleName }}"
        }
        this.node = choice
    }

    inline fun <reified T : Any> our(configure: OurChoice<T>.() -> Unit) {
        val tagClass = T::class
        val annotation = tagClass.findAnnotation<Alpha>()
        requireNotNull(annotation) {
            "Class '${tagClass.simpleName}' must have '${Alpha::class.simpleName}' annotation"
        }
        require(tagClass.isSealed) {
            "Class '${tagClass.simpleName}' must be sealed"
        }
        val choice = OurChoice<T>()
        choice.configure()
        this.node = choice
    }

    fun end() {
        this.node = End()
    }
}

class Tree : TreeNode()

data class WeFromThem<T>(val name: String) : TreeNode()
data class TheyFromUs<T>(val name: String) : TreeNode()

class End : TreeNode()

class TheirChoice<T : Any> : TreeNode() {
    val branches = mutableMapOf<KClass<out T>, WeFromThem<out T>>()

    fun choice(messageClass: KClass<out T>, configure: WeFromThem<out T>.() -> Unit) {
        val annotation = messageClass.findAnnotation<State>()
        requireNotNull(annotation) {
            "Class '${messageClass.simpleName}' must have '${State::class.simpleName}' annotation"
        }
        val branch = WeFromThem<T>(annotation.name)
        branch.configure()
        branches[messageClass] = branch
    }
}

class OurChoice<T : Any> : TreeNode() {
    val branches = mutableMapOf<KClass<out T>, TheyFromUs<out T>>()

    fun choice(messageClass: KClass<out T>, configure: TheyFromUs<out T>.() -> Unit) {
        val annotation = messageClass.findAnnotation<State>()
        requireNotNull(annotation) {
            "Class '${messageClass.simpleName}' must have '${State::class.simpleName}' annotation"
        }
        val branch = TheyFromUs<T>(annotation.name)
        branch.configure()
        branches[messageClass] = branch
    }
}

@DslMarker
annotation class TreeMarker
