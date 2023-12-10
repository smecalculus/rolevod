package smecalculus.bezmen.storage;

import java.util.Optional;
import java.util.UUID;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import smecalculus.bezmen.core.SepulkaStateDm;
import smecalculus.bezmen.storage.springdata.SepulkaRepository;

@RequiredArgsConstructor
public class SepulkaDaoSpringData implements SepulkaDao {

    @NonNull
    private SepulkaStateMapper mapper;

    @NonNull
    private SepulkaRepository repository;

    @Override
    public SepulkaStateDm.AggregateRoot add(@NonNull SepulkaStateDm.AggregateRoot state) {
        var stateEdge = repository.save(mapper.toEdge(state));
        return mapper.toDomain(stateEdge);
    }

    @Override
    public Optional<SepulkaStateDm.Existence> getBy(@NonNull String externalId) {
        return repository.findByExternalId(externalId).map(mapper::toDomain);
    }

    @Override
    public Optional<SepulkaStateDm.Preview> getBy(@NonNull UUID internalId) {
        return repository.findByInternalId(internalId).map(mapper::toDomain);
    }

    @Override
    public void updateBy(UUID internalId, SepulkaStateDm.Touch state) {
        var stateEdge = mapper.toEdge(state);
        var matchedCount = repository.updateBy(internalId, stateEdge);
        if (matchedCount == 0) {
            throw new ContentionException();
        }
    }
}
